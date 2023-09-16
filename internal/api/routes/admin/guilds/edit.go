package guilds

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asianchinaboi/backendserver/internal/api/middleware"
	"github.com/asianchinaboi/backendserver/internal/config"
	"github.com/asianchinaboi/backendserver/internal/db"
	"github.com/asianchinaboi/backendserver/internal/errors"
	"github.com/asianchinaboi/backendserver/internal/events"
	"github.com/asianchinaboi/backendserver/internal/files"
	"github.com/asianchinaboi/backendserver/internal/logger"
	"github.com/asianchinaboi/backendserver/internal/session"
	"github.com/asianchinaboi/backendserver/internal/uid"
	"github.com/asianchinaboi/backendserver/internal/wsclient"
	"github.com/gin-gonic/gin"
)

type editGuildBody struct {
	SaveChat *bool   `json:"saveChat"`
	Name     *string `json:"name"`
	OwnerId  *int64  `json:"ownerId"`
}

func Edit(c *gin.Context) {
	user := c.MustGet(middleware.User).(*session.Session)
	if user == nil {
		errors.SendErrorResponse(c, errors.ErrSessionDidntPass, errors.StatusInternalError)
		return
	}
	if !user.Perms.Admin && !user.Perms.Guilds.Edit {
		errors.SendErrorResponse(c, errors.ErrNotAuthorised, errors.StatusNotAuthorised)
		return
	}
	guildId := c.Param("guildId")
	if match, err := regexp.MatchString("^[0-9]+$", guildId); err != nil {
		errors.SendErrorResponse(c, err, errors.StatusInternalError)
		return
	} else if !match {
		errors.SendErrorResponse(c, errors.ErrRouteParamInvalid, errors.StatusRouteParamInvalid)
		return
	}

	var newSettings editGuildBody
	var imageHeader *multipart.FileHeader
	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		var err error
		if imageHeader, err = c.FormFile("image"); err != nil && err != http.ErrMissingFile {
			errors.SendErrorResponse(c, err, errors.StatusBadRequest)
			return
		}
		jsonData := c.PostForm("body")
		if err := json.Unmarshal([]byte(jsonData), &newSettings); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusBadRequest)
			return
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		if err := c.ShouldBindJSON(&newSettings); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusBadRequest)
			return
		}
	} else {
		errors.SendErrorResponse(c, errors.ErrNotSupportedContentType, errors.StatusBadRequest)
		return
	}

	if newSettings.SaveChat == nil && newSettings.Name == nil && newSettings.OwnerId == nil && imageHeader == nil {
		errors.SendErrorResponse(c, errors.ErrAllFieldsEmpty, errors.StatusAllFieldsEmpty)
		return
	}

	successful := false

	var exists bool
	if err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM guilds WHERE id = $1)", guildId).Scan(&exists); err != nil {
		errors.SendErrorResponse(c, err, errors.StatusInternalError)
		return
	}
	if !exists {
		errors.SendErrorResponse(c, errors.ErrGuildNotExist, errors.StatusGuildNotExist)
		return
	}

	bodyRes := events.Guild{}
	intGuildId, err := strconv.ParseInt(guildId, 10, 64)
	if err != nil {
		errors.SendErrorResponse(c, err, errors.StatusInternalError)
		return
	}

	//BEGIN TRANSACTION
	ctx := context.Background()
	tx, err := db.Db.BeginTx(ctx, nil)
	if err != nil {
		errors.SendErrorResponse(c, err, errors.StatusInternalError)
		return
	}
	defer tx.Rollback() //rollback changes if failed

	bodyRes.GuildId = intGuildId

	bodySettingsRes := events.Guild{}

	bodySettingsRes.GuildId = intGuildId

	if newSettings.SaveChat != nil {

		if _, err = tx.ExecContext(ctx, "UPDATE guilds SET save_chat=$1 WHERE id=$2", *newSettings.SaveChat, guildId); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodySettingsRes.SaveChat = newSettings.SaveChat
	} else {
		var saveChat bool
		if err := db.Db.QueryRow("SELECT save_chat FROM guilds WHERE id = $1", guildId).Scan(&saveChat); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodySettingsRes.SaveChat = &saveChat
	}
	if newSettings.Name != nil {
		if _, err = tx.ExecContext(ctx, "UPDATE guilds SET name=$1 WHERE id=$2", *newSettings.Name, guildId); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodyRes.Name = *newSettings.Name
	} else {
		var name string
		if err := db.Db.QueryRow("SELECT name FROM guilds WHERE id = $1", guildId).Scan(&name); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodyRes.Name = name
	}

	if imageHeader != nil {
		//remove old image

		var oldImageId int64
		if err := tx.QueryRowContext(ctx, "DELETE FROM files WHERE guild_id = $1 RETURNING id", guildId).Scan(&oldImageId); err != nil && err != sql.ErrNoRows {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}

		imageId := uid.Snowflake.Generate().Int64()
		filename := imageHeader.Filename
		fileType := filepath.Ext(filename)
		imageCreated := time.Now().Unix()
		image, err := imageHeader.Open()
		if err != nil {
			errors.SendErrorResponse(c, err, errors.StatusBadRequest)
			return
		}
		defer image.Close()

		fileBytes, err := io.ReadAll(image)
		if err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}

		if valid := files.ValidateImage(fileBytes, fileType); !valid {
			errors.SendErrorResponse(c, errors.ErrFileInvalid, errors.StatusFileInvalid)
			return
		}

		filesize := len(fileBytes)

		if filesize > config.Config.Server.MaxFileSize {
			errors.SendErrorResponse(c, errors.ErrFileTooLarge, errors.StatusFileTooLarge)
			return
		} else if !(filesize >= 0) {
			errors.SendErrorResponse(c, errors.ErrFileNoBytes, errors.StatusFileNoBytes)
			return
		}

		compressedBuffer, err := files.Compress(fileBytes, filesize)
		if err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}

		outFile, err := os.Create(fmt.Sprintf("uploads/guild/%d.lz4", imageId))
		if err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		defer func() { //defer just in case something went wrong
			if !successful {
				if err := os.Remove(fmt.Sprintf("uploads/guild/%d.lz4", imageId)); err != nil {
					logger.Warn.Printf("failed to remove file: %v\n", err)
				}
			} else {
				if oldImageId != 0 {
					if err := os.Remove(fmt.Sprintf("uploads/guild/%d.lz4", oldImageId)); err != nil {
						logger.Warn.Printf("failed to remove file: %v\n", err)
					}
				}
			}
		}()
		defer outFile.Close()

		if _, err = outFile.Write(compressedBuffer); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}

		if _, err = tx.ExecContext(ctx, "INSERT INTO files (id, guild_id, filename, created, temp, filesize, entity_type) VALUES ($1, $2, $3, $4, $5, $6, 'guild')", imageId, guildId, filename, imageCreated, false, filesize); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodyRes.ImageId = imageId
	} else {
		var imageId int64
		if err := db.Db.QueryRow("SELECT id FROM files WHERE guild_id=$1", guildId).Scan(&imageId); err != nil && err != sql.ErrNoRows {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		} else if err == sql.ErrNoRows {
			bodyRes.ImageId = -1
		} else {
			bodyRes.ImageId = imageId
		}
	}

	if newSettings.OwnerId != nil {
		if _, err = tx.ExecContext(ctx, "UPDATE userguilds SET owner=false WHERE guild_id=$1 AND owner = true", guildId); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO userguilds(user_id, guild_id, owner) VALUES ($1, $2, true) ON CONFLICT (guild_id, user_id) DO UPDATE SET owner = true", newSettings.OwnerId, guildId); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
		bodyRes.OwnerId = *newSettings.OwnerId
	} else {
		if err := db.Db.QueryRow("SELECT user_id FROM userguilds WHERE guild_id=$1 AND owner = true", guildId).Scan(&bodyRes.OwnerId); err != nil {
			errors.SendErrorResponse(c, err, errors.StatusInternalError)
			return
		}
	}

	if err := tx.Commit(); err != nil { //commits the transaction
		errors.SendErrorResponse(c, err, errors.StatusInternalError)
		return
	}
	successful = true

	guildRes := wsclient.DataFrame{
		Op:    wsclient.TYPE_DISPATCH,
		Data:  bodyRes,
		Event: events.UPDATE_GUILD,
	}
	wsclient.Pools.BroadcastGuild(intGuildId, guildRes)

	c.Status(http.StatusNoContent)
}

// mautrix-telegram - A Matrix-Telegram puppeting bridge.
// Copyright (C) 2026 Sumner Evans
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-telegram/pkg/gotd/tg"
)

func TestMakeViewOnce(t *testing.T) {
	photo := &tg.InputMediaUploadedPhoto{}
	assert.NoError(t, makeViewOnce(photo, event.MsgImage))
	ttl, ok := photo.GetTTLSeconds()
	assert.True(t, ok)
	assert.Equal(t, viewOnceTTL, ttl)

	video := &tg.InputMediaUploadedDocument{MimeType: "video/mp4"}
	assert.NoError(t, makeViewOnce(video, event.MsgVideo))
	ttl, ok = video.GetTTLSeconds()
	assert.True(t, ok)
	assert.Equal(t, viewOnceTTL, ttl)

	// Images that had to be sent as documents can't be view-once, and neither can non-video files
	doc := &tg.InputMediaUploadedDocument{MimeType: "image/png"}
	assert.ErrorIs(t, makeViewOnce(doc, event.MsgImage), bridgev2.ErrViewOnceNotAllowed)
	assert.ErrorIs(t, makeViewOnce(&tg.InputMediaUploadedDocument{MimeType: "application/octet-stream"}, event.MsgVideo), bridgev2.ErrViewOnceNotAllowed)
	_, ok = doc.GetTTLSeconds()
	assert.False(t, ok)
}

func TestViewOnceFileCaps(t *testing.T) {
	for _, caps := range []event.FileFeatureMap{dmFileCaps, premiumDMFileCaps} {
		assert.True(t, caps[event.MsgImage].ViewOnce)
		assert.True(t, caps[event.MsgVideo].ViewOnce)
		assert.False(t, caps[event.MsgFile].ViewOnce)
	}
	for _, caps := range []event.FileFeatureMap{fileCaps, premiumFileCaps} {
		assert.False(t, caps[event.MsgImage].ViewOnce)
		assert.False(t, caps[event.MsgVideo].ViewOnce)
	}
}

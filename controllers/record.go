package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nelwhix/Attendlog/models"
	"github.com/Nelwhix/Attendlog/requests"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
	"log"
	"net/http"
	"strings"
	"time"
)

func (c *Controller) CreateNewRecord(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		flashMessage := map[string]string{
			"type":    "error",
			"message": "error parsing form",
		}

		ctx := context.WithValue(r.Context(), "flashMessage", flashMessage)
		r = r.WithContext(ctx)
		c.RenderLinkForm(w, r)
		log.Printf("error parsing form: %v", err.Error())
		return
	}

	r.PostForm.Del("gorilla.csrf.Token")

	var nRecordReq requests.Record
	err = decoder.Decode(&nRecordReq, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		flashMessage := map[string]string{
			"type":    "error",
			"message": err.Error(),
		}
		ctx := context.WithValue(r.Context(), "flashMessage", flashMessage)
		r = r.WithContext(ctx)
		c.RenderLinkForm(w, r)
		log.Printf("error decoding form: %v", err.Error())
		return
	}

	pathParams := mux.Vars(r)
	nRecord := &models.Record{
		ID:        ulid.Make().String(),
		FirstName: nRecordReq.FirstName,
		LastName:  nRecordReq.LastName,
		Email:     nRecordReq.Email,
		LinkID:    pathParams["id"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := c.db.Create(&nRecord)
	if result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error inserting record: %v", result.Error.Error())
		return
	}

	flashMessage := map[string]string{
		"type":    "success",
		"message": "details submitted! enjoy the event",
	}
	ctx := context.WithValue(r.Context(), "flashMessage", flashMessage)
	r = r.WithContext(ctx)
	c.RenderLinkForm(w, r)
}

func (c *Controller) GetRecords(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server Sent Events not supported by client", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/event-stream")
	recordsCh := make(chan []models.Record)
	pathParams := mux.Vars(r)
	go c.getNewRecords(r.Context(), recordsCh, pathParams["id"])

	for records := range recordsCh {
		event, err := c.formatSSE(records)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			break
		}

		_, err = fmt.Fprint(w, event)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			break
		}

		flusher.Flush()
	}
}

func (c *Controller) formatSSE(records []models.Record) (string, error) {
	recordsJson, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	m := map[string]interface{}{
		"data": map[string]string{
			"message": "Get Records.",
			"records": string(recordsJson),
		},
	}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err = encoder.Encode(m)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}

func (c *Controller) getNewRecords(ctx context.Context, recordsChan chan<- []models.Record, linkID string) {
	ticker := time.NewTicker(time.Minute)
	var count1 int64
	var count2 int64
	c.db.Model(&models.Record{}).Where("link_id = ?", linkID).Count(&count1)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		case <-ticker.C:
			c.db.Model(&models.Record{}).Where("link_id = ?", linkID).Count(&count2)
			if count2 > count1 {
				var records []models.Record
				c.db.Where("link_id = ?", linkID).Limit(50).Find(&records)
				recordsChan <- records
			}
		}
	}

	ticker.Stop()
	close(recordsChan)
}

package cmd

import (
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/steipete/gogcli/internal/ui"
)

type attachmentInfo struct {
	Filename     string
	Size         int64
	MimeType     string
	AttachmentID string
}

type attachmentOutput struct {
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	SizeHuman    string `json:"sizeHuman"`
	MimeType     string `json:"mimeType"`
	AttachmentID string `json:"attachmentId"`
}

type attachmentDownloadOutput struct {
	MessageID string `json:"messageId"`
	attachmentOutput
	Path   string `json:"path,omitempty"`
	Cached bool   `json:"cached,omitempty"`
}

func attachmentOutputFromInfo(a attachmentInfo) attachmentOutput {
	return attachmentOutput{
		Filename:     a.Filename,
		Size:         a.Size,
		SizeHuman:    formatBytes(a.Size),
		MimeType:     a.MimeType,
		AttachmentID: a.AttachmentID,
	}
}

func attachmentOutputs(attachments []attachmentInfo) []attachmentOutput {
	if len(attachments) == 0 {
		return nil
	}
	out := make([]attachmentOutput, len(attachments))
	for i, a := range attachments {
		out[i] = attachmentOutputFromInfo(a)
	}
	return out
}

func attachmentOutputsFromDownloads(attachments []attachmentDownloadOutput) []attachmentOutput {
	if len(attachments) == 0 {
		return nil
	}
	out := make([]attachmentOutput, len(attachments))
	for i, a := range attachments {
		out[i] = a.attachmentOutput
	}
	return out
}

func attachmentLine(a attachmentOutput) string {
	return fmt.Sprintf("attachment\t%s\t%s\t%s\t%s", a.Filename, a.SizeHuman, a.MimeType, a.AttachmentID)
}

func printAttachmentLines(p *ui.Printer, attachments []attachmentOutput) {
	for _, a := range attachments {
		p.Println(attachmentLine(a))
	}
}

func printAttachmentSection(p *ui.Printer, attachments []attachmentInfo) {
	out := attachmentOutputs(attachments)
	if len(out) == 0 {
		return
	}
	p.Println("Attachments:")
	printAttachmentLines(p, out)
	p.Println("")
}

func collectAttachments(p *gmail.MessagePart) []attachmentInfo {
	if p == nil {
		return nil
	}
	var out []attachmentInfo
	if p.Body != nil && p.Body.AttachmentId != "" {
		filename := p.Filename
		if strings.TrimSpace(filename) == "" {
			filename = "attachment"
		}
		out = append(out, attachmentInfo{
			Filename:     filename,
			Size:         p.Body.Size,
			MimeType:     p.MimeType,
			AttachmentID: p.Body.AttachmentId,
		})
	}
	for _, part := range p.Parts {
		out = append(out, collectAttachments(part)...)
	}
	return out
}

// formatBytes formats bytes into human-readable format.
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

package cmd

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

func buildEventDateTime(value string, allDay bool) *calendar.EventDateTime {
	value = strings.TrimSpace(value)
	if allDay {
		return &calendar.EventDateTime{Date: value}
	}
	return &calendar.EventDateTime{DateTime: value}
}

func buildConferenceData(withMeet bool) *calendar.ConferenceData {
	if !withMeet {
		return nil
	}
	return &calendar.ConferenceData{
		CreateRequest: &calendar.CreateConferenceRequest{
			RequestId: fmt.Sprintf("gogcli-%d", time.Now().UnixNano()),
			ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
				Type: "hangoutsMeet",
			},
		},
	}
}

func buildRecurrence(rules []string) []string {
	if len(rules) == 0 {
		return nil
	}
	out := make([]string, 0, len(rules))
	for _, r := range rules {
		r = strings.TrimSpace(r)
		if r != "" {
			out = append(out, r)
		}
	}
	return out
}

func buildAttachments(urls []string) []*calendar.EventAttachment {
	if len(urls) == 0 {
		return nil
	}
	out := make([]*calendar.EventAttachment, 0, len(urls))
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			out = append(out, &calendar.EventAttachment{FileUrl: u})
		}
	}
	return out
}

func buildExtendedProperties(privateProps, sharedProps []string) *calendar.EventExtendedProperties {
	if len(privateProps) == 0 && len(sharedProps) == 0 {
		return nil
	}
	props := &calendar.EventExtendedProperties{}

	if len(privateProps) > 0 {
		props.Private = make(map[string]string)
		for _, p := range privateProps {
			if k, v, ok := strings.Cut(p, "="); ok {
				props.Private[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		}
	}

	if len(sharedProps) > 0 {
		props.Shared = make(map[string]string)
		for _, p := range sharedProps {
			if k, v, ok := strings.Cut(p, "="); ok {
				props.Shared[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		}
	}

	return props
}

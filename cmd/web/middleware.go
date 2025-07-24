package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type ContextKey string

const ContextUserKey ContextKey = "user_ip"

func (app *State) IpFromContext(ctx context.Context) string {
	return ctx.Value(ContextUserKey).(string)
}

func (app *State) AddIpToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := GetIP(r)
		fmt.Println(ip)
		if err != nil {
			ip = "unknown"
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, ip)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}

	user_ip := net.ParseIP(ip)
	if user_ip == nil {
		return "", fmt.Errorf("user ip: %q is not IP:port", r.RemoteAddr)
	}

	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return "unknown", err
}

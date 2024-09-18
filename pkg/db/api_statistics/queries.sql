-- name: InsertApiStatistic :exec 
INSERT INTO public.api_statistics (url, status_code, http_method, size_bytes, client_ip, elapsed_time_seconds) 
VALUES($1, $2, $3, $4, $5, $6);

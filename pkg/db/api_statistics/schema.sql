CREATE TABLE api_statistics (
	timestamp timestamp NOT NULL DEFAULT NOW(),
	url TEXT NOT NULL,
	status_code int NOT NULL,
	http_method TEXT NOT NULL,
	size_bytes bigint NOT NULL,
	client_ip TEXT NOT NULL,
	elapsed_time_seconds float8
);
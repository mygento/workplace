package nginx

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	RuntimeDir string
	Root       string
	Host       string
	Port       int
	FPMSocket  string
}

func WriteConfig(cfg Config) (string, error) {
	if err := os.MkdirAll(
		cfg.RuntimeDir,
		0755,
	); err != nil {
		return "", fmt.Errorf(
			"create Nginx runtime directory: %w",
			err,
		)
	}

	configPath := filepath.Join(
		cfg.RuntimeDir,
		"nginx.conf",
	)

	pidPath := filepath.Join(
		cfg.RuntimeDir,
		"nginx.pid",
	)

	accessLog := filepath.Join(
		cfg.RuntimeDir,
		"access.log",
	)

	errorLog := filepath.Join(
		cfg.RuntimeDir,
		"error.log",
	)

	content := fmt.Sprintf(`worker_processes 1;

pid %s;

error_log %s;

events {
    worker_connections 1024;
}

http {

    access_log %s;

    default_type application/octet-stream;

    sendfile on;

    types {
        text/html html;
        text/css css;
        text/plain txt;

        application/javascript js;
        application/json json;
        application/xml xml;
        application/pdf pdf;

        image/png png;
        image/jpeg jpg jpeg;
        image/gif gif;
        image/svg+xml svg;
        image/webp webp;
        image/x-icon ico;

        font/woff woff;
        font/woff2 woff2;
        font/ttf ttf;
    }

    server {

        listen %s:%d;

        server_name localhost;

        root %s;

        index index.php index.html;


        location / {
            try_files $uri $uri/ /index.php?$query_string;
        }


        location ~ \.php$ {

            try_files $uri =404;

            include fastcgi_params;

            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            fastcgi_param DOCUMENT_ROOT $document_root;

            fastcgi_param HTTP_PROXY "";

            fastcgi_pass unix:%s;
        }


        location ~ /\. {
            deny all;
        }
    }
}
`,
		pidPath,
		errorLog,
		accessLog,
		cfg.Host,
		cfg.Port,
		cfg.Root,
		cfg.FPMSocket,
	)

	if err := os.WriteFile(
		configPath,
		[]byte(content),
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"write Nginx config: %w",
			err,
		)
	}

	fastcgiParamsPath := filepath.Join(
		cfg.RuntimeDir,
		"fastcgi_params",
	)
	fastcgiParams := `fastcgi_param QUERY_STRING       $query_string;
fastcgi_param REQUEST_METHOD     $request_method;
fastcgi_param CONTENT_TYPE       $content_type;
fastcgi_param CONTENT_LENGTH     $content_length;

fastcgi_param SCRIPT_NAME        $fastcgi_script_name;
fastcgi_param REQUEST_URI        $request_uri;
fastcgi_param DOCUMENT_URI       $document_uri;
fastcgi_param DOCUMENT_ROOT      $document_root;
fastcgi_param SERVER_PROTOCOL    $server_protocol;
fastcgi_param REQUEST_SCHEME     $scheme;
fastcgi_param HTTPS              $https if_not_empty;

fastcgi_param GATEWAY_INTERFACE  CGI/1.1;
fastcgi_param SERVER_SOFTWARE    nginx/$nginx_version;

fastcgi_param REMOTE_ADDR        $remote_addr;
fastcgi_param REMOTE_PORT        $remote_port;
fastcgi_param SERVER_ADDR        $server_addr;
fastcgi_param SERVER_PORT        $server_port;
fastcgi_param SERVER_NAME        $server_name;

fastcgi_param REDIRECT_STATUS    200;
`
	if err := os.WriteFile(
		fastcgiParamsPath,
		[]byte(fastcgiParams),
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"write fastcgi parameters: %w",
			err,
		)
	}

	return configPath, nil
}

# CDN S3 Backend Docs

### Minimum Requirement (LOM)
<b>

- 2 GB Memory
- 1 Core CPU
- 1 GB Storage
- Docker
- Docker-Compose Plugin
- Certbot
- At least one IPv4 pointing to a domain

</b>

### Deployment Instructions
- Clone the project from CVS:
  ```shell
  git clone <repo> && cd <app_dir>
  ```
- Copy <code>.env.example</code> to <code>.env</code> And modify it as you wish.
  ```shell
  cp .env.example .env
  ```
- Obtain necessary SSL Certificates with Certbot:
  ```shell
  sudo certbot certonly --manual --preferred-challenges=dns -d *.APP_HOST -d $APP_HOST
  ```
- Manage TLS certs:
  ```shell
  mkdir -p build/certs/{app,nginx,minio,redis}
  cp /etc/letsencrypt/live/$APP_HOST/* build/certs/app
  cp /etc/letsencrypt/live/$APP_HOST/* build/certs/nginx
  cp /etc/letsencrypt/live/$APP_HOST/* build/certs/redis
  cp /etc/letsencrypt/live/$APP_HOST/privkey.pem build/certs/minio/private.key
  cp /etc/letsencrypt/live/$APP_HOST/fullchain.pem build/certs/minio/public.crt
  chmod 755 build/certs/redis/*
  ```
- Bootstrap the application via:
  ```shell
  docker compose up -d
  ```
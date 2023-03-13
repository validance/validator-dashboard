FROM golang:1.19.7-alpine

# Go 프로젝트 코드 복사
COPY . /app

# 작업 디렉토리로 이동
WORKDIR /app/app

# 필요한 종속성 설치
RUN go mod download
# timezone 설정
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Seoul /etc/localtime  \
    && echo "Asia/Seoul" > /etc/timezone

# 프로젝트 빌드
RUN go build -o ../dashboard

WORKDIR /app
# 실행
CMD ["./dashboard"]
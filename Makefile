BUILDDIR=${CURDIR}/build
BINNAME=users
APP_ENTRYPOINT=./cmd/users/users.go
JET_GENERATE_PATH=./internal/pkg/jet

clean:
	rm -rf ${BUILDDIR}

bindir:
	mkdir -p ${BUILDDIR}

build: bindir
	go build -o ${BUILDDIR}/${BINNAME} ${APP_ENTRYPOINT}

start:
	${BUILDDIR}/${BINNAME}


# generate targets
generate-all: generate-jet

generate-jet:
	jet -dsn=postgres://postgres:password@localhost:5432/pechat-users?sslmode=disable -schema=pechat-users -path=${JET_GENERATE_PATH}


# tools installation targets
install-all: install-jet-generator

install-jet-generator:
	# TODO зафиксировать версию
	go install github.com/go-jet/jet/v2/cmd/jet@latest
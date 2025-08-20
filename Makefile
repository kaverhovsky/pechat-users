BUILDDIR=${CURDIR}/build
BINDIR=${CURDIR}/bin
MIGRATIONDIR=${CURDIR}/migrations
BINNAME=users
APP_ENTRYPOINT=./cmd/users/users.go
JET_GENERATE_PATH=./internal/pkg/jet

clean:
	rm -rf ${BUILDDIR}

bindir:
	mkdir -p ${BINDIR}

builddir:
	mkdir -p ${BUILDDIR}

build: bindir
	go build -o ${BUILDDIR}/${BINNAME} ${APP_ENTRYPOINT}

start:
	${BUILDDIR}/${BINNAME}


# generate targets
generate-all: generate-jet

generate-jet:
	${BINDIR}/jet \
	-dsn=postgres://postgres:password@localhost:5432/pechat_users?sslmode=disable \
	-path=${JET_GENERATE_PATH} \
	-ignore-tables goose_db_version, goose_db_version_id_seq


# postgres migrations
up-goose-migrations:
	${BINDIR}/goose postgres "postgres://postgres:password@localhost:5432/pechat_users?sslmode=disable" -dir ${MIGRATIONDIR} up

down-goose-migrations:
	${BINDIR}/goose postgres "postgres://postgres:password@localhost:5432/pechat_users?sslmode=disable" -dir ${MIGRATIONDIR} down

# tools installation targets
install-all: bindir install-jet-generator install-goose

install-jet-generator:
	$(info Installing jet binary into [$(BINDIR)]...)
	GOBIN=$(BINDIR) go install github.com/go-jet/jet/v2/cmd/jet@v2.13.0

install-goose:
	$(info Installing goose binary into [$(BINDIR)]...)
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
Yumster
=======

`yumster` is a fork of https://github.com/FINRAOS/yum-nginx-api for simple yum repository management.

## How to build yumster (Docker image)

    git clone https://github.com/hyung-hwan/yumster.git
    cd yumster
    make

## How to Install yumster (Binary)

    make build

## Configuration File `yumster.yml`

**Configuration file can be JSON, TOML, YAML, HCL, or Java properties**

    # createrepo workers, default is 2
    createrepo_workers:
    # http max content upload, default is 10000000 <- 10MB
    max_content_length:
    # yum repo directory, default is ./
    upload_dir:
    # port to run http server, default is 8080
    port:
    # max retries to retry failed createrepo, default is 3
    max_retries:

## Run it as a container

    docker run -d -p 8080:8080 --name yumster yumster:latest

## API Usage

**Post binary RPM to the my-repo repository using the API endpoint:**

    curl -F file=@abc.rpm http://localhost:8080/api/upload/my-repo

**Get repomd.xml of the my-repo repository:**

    curl http://localhost:8080/repo/my-repo/repodata/repomd.xml

**Get the uploaded file

    curl -o /tmp/abc.rpm http://localhost:8080/repo/my-repo/abc.rpm

**Health check API endpoint**

    curl http://localhost/api/health

## License Type

yumster project is follows the license of `yum-nginx-api`.

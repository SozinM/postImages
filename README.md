# Post Images
[![Build Status](https://travis-ci.org/SozinM/postImages.svg?branch=master)](https://travis-ci.org/SozinM/postImages)

Post image is rest service for (only) uploading images and creating thumbnails.
Entry point for service is "/images"
Json format for request is {"url":[list of urls]}, for example:
    {
        "url": [
            "https://blog.golang.org/go-image-package_image-package-01.png",
            "1",
            "https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png"
        ]
    }
Service answer contains urls and success status, for example:
    [
        {
            "url": "https://blog.golang.org/go-image-package_image-package-01.png",
            "success": true
        },
        {
            "url": "1",
            "success": false
        },
        {
            "url": "https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png",
            "success": true
        }
    ]

# How to build:
    cd ~
    git clone https://github.com/SozinM/postImages
    cd postImages
    go test -v -race
    go build

# Run inside container
    cd ~
    git clone https://github.com/SozinM/postImages
    cd postImages
    docker-compose up

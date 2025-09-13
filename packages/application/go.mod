module govel/application

go 1.25

require (
    // Add your dependencies here
)


replace (
    // Map the main module to the src directory
    govel/application => ./src

    // Inter-package dependencies
    govel/config => ../config
    govel/container => ../container
    govel/cookie => ../cookie
    govel/encryption => ../encryption
    govel/hashing => ../hashing
    govel/logger => ../logger
    govel/pipeline => ../pipeline
    govel/support => ../support
    govel/types => ../types

)

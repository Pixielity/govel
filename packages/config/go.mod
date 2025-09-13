module govel/config

go 1.25


replace (
    // Map the main module to the src directory
    govel/config => ./src

    // Inter-package dependencies
    govel/application => ../application
    govel/container => ../container
    govel/cookie => ../cookie
    govel/encryption => ../encryption
    govel/hashing => ../hashing
    govel/logger => ../logger
    govel/pipeline => ../pipeline
    govel/support => ../support
    govel/types => ../types

)

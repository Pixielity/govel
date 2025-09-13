module govel/logger

go 1.25

require (
    // Add your dependencies here
)


replace (
    // Map the main module to the src directory
    govel/logger => ./src

    // Inter-package dependencies
    govel/application => ../application
    govel/config => ../config
    govel/container => ../container
    govel/cookie => ../cookie
    govel/encryption => ../encryption
    govel/hashing => ../hashing
    govel/pipeline => ../pipeline
    govel/support => ../support
    govel/types => ../types

)

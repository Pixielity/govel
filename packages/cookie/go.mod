module govel/cookie

go 1.25

require (
    // Add your dependencies here
)


replace (
    // Map the main module to the src directory
    govel/cookie => ./src

    // Inter-package dependencies
    govel/application => ../application
    govel/config => ../config
    govel/container => ../container
    govel/encryption => ../encryption
    govel/hashing => ../hashing
    govel/logger => ../logger
    govel/pipeline => ../pipeline
    govel/support => ../support
    govel/types => ../types

)

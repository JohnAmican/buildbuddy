load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(default_visibility = ["//enterprise:__subpackages__"])

go_library(
    name = "redis_client",
    srcs = ["redis_client.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy/enterprise/server/backends/redis_client",
    deps = [
        "//enterprise/server/util/redisutil",
        "//server/real_environment",
        "//server/remote_execution/config",
        "//server/util/flag",
        "//server/util/status",
    ],
)

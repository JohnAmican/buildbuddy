load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(default_visibility = ["//enterprise:__subpackages__"])

go_library(
    name = "invocation_search_service",
    srcs = ["invocation_search_service.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy/enterprise/server/invocation_search_service",
    deps = [
        "//proto:invocation_go_proto",
        "//proto:invocation_status_go_proto",
        "//server/build_event_protocol/build_event_handler",
        "//server/build_event_protocol/invocation_format",
        "//server/environment",
        "//server/interfaces",
        "//server/tables",
        "//server/util/alert",
        "//server/util/blocklist",
        "//server/util/clickhouse/schema",
        "//server/util/db",
        "//server/util/filter",
        "//server/util/git",
        "//server/util/perms",
        "//server/util/query_builder",
        "//server/util/status",
        "//server/util/uuid",
    ],
)

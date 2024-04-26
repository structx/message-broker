
job "mora" {

    datacenters = ["dc1"]
    type = "service"

    group "blockchain" {
        count = 1

        network {
            mode = "bridge"

            port "http" {}

            port "rpc" {}

            port "metrics" {
                to = 2112
            }
        }

        service {
            name = "mora-http-api" 
            port = "http"

            connect {
                sidecar_service {}
            }

            check {
                type = "http"
                path = "/health"
                interval = "10s"
                timeout = "3s"
            }
        }

        service {
            name = "mora-metrics"
            port = "metrics"

            tags = [
                "metrics"
            ]

            provider = "consul"
        }

        service {
            name = "mora-rpc"
            port = "rpc"

            provider = "consul"

            connect {
                sidecar_service {}
            }
        }

        task "server" {
            driver = "docker"

            config {
                image = "trevatk/mora:v0.0.1"
                ports = [ "http", "rpc", "metrics" ]
            }

            env {
                DSERVICE_CONFIG = "${NOMAD_TASK_DIR}/config.hcl"
            }

            template {
                destination = "local/config.hcl"
                change_mode = "signal"
                change_signal = "SIGTERM"
                data = <<EOH
server {
    bind_addr = "0.0.0.0"

    ports {
        http = {{ env "NOMAD_PORT_http"}}
        grpc = {{ env "NOMAD_PORT_rpc" }}
    }
}

raft {
    bootstrap = true
    local_id = "1234567789012345456"
    base_dir = "/opt/mora/raft"
}

logger {
    log_path = "/var/log/mora/node.log"
    log_level = "DEBUG"
    raft_log_path = "/var/log/mora"
}

chain {
    base_dir = ""
}

message_broker {
    server_addr = ""
}
                EOH
            }

        }
    }
}
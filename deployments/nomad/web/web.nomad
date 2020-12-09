job "web"{
  region = "global"
  datacenters = ["dc1"]
  type = "service"

  update {
    stagger = "30s"
    max_parallel = 1
  }

  group "web"{
      count = 1

      task "web"{
          driver = "docker"
          image = "cmd/web:latest"

          port_map {
              http = $ADDR
          }

        env{
            ADDR = "8080"
            REDIS_ADDR = "${NOMAD_IP}:6379"
            REDIS_MESSAGES_CHANNEL = "messages"
            REDIS_RESPONSES_CHANNEL = "responses"
        }

      }
  }

  # Service block needed for consul health check and service monitoring
  service {
      port "http"

      check {
          type = "http"
          path = "/"
          interval = "10sec"
          timeout = "3sec"

      }
  }

  resources {
      cpu = 128
      mem = 128
      network{
          mbits = 100
          port "http" {}
      }
  }
}
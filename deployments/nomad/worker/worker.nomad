job "worker"{
  region = "global"
  datacenters = ["dc1"]
  type = "service"

  update {
    stagger = "30s"
    max_parallel = 1
  }

  group "worker"{
      count = 1

      task "web"{
          driver = "docker"
          image = "cmd/worker:latest"
          
          env{
            REDIS_ADDR = "${NOMAD_IP}:6379"
            REDIS_CHANNEL = "messages"
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
      }
  }
}
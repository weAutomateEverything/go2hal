provider "kubernetes" {
  config_context_auth_info = "${var.k8s_auth_info}"
  config_context_cluster = "${var.k8s_context_cluster}"
}


resource "kubernetes_config_map" "mongo" {
  "metadata" {
    name = "mongo"
    namespace = "go2hal"
  }
  data {
    MONGO_SERVERS = "${var.mongo_servers}"
    MONGO_USER = "${var.mongo_user}"
    MONGO_REPLICA_SET = "${var.mongo_replica_set}"
    MONGO_AUTH_SOURCE = "${var.mongo_auth_source}"
    MONGO_SSL = "${var.mongo_ssl}"
  }
}

resource "kubernetes_secret" "mongo" {
  "metadata" {
    name = "mongo"
    namespace = "go2hal"
  }
  data {
    MONGO_PASSWORD = "${var.mongo_password}"
  }
}

resource "kubernetes_secret" "go2hal" {
  "metadata" {
    name = "go2hal"
    namespace = "go2hal"
  }

  data {
    HTTP_PROXY = "${var.http_proxy}"
    HTTPS_PROXY = "${var.https_proxy}"
    BOT_KEY =  "${var.bot_key}"
  }
}

resource "kubernetes_pod" "go2hal" {
  "metadata" {
    labels {
      app = "hal"
    }
    name = "hal"
    namespace = "go2hal"
  }
  count = 1
  "spec" {
    container {
      name = "go2hal"
      image = "weautomateeverything/go2hal:latest"
      port {
        container_port = 8000
        protocol = "TCP"
      }
      port {
        container_port = 8080
        protocol = "TCP"
      }
      env {
        name = "MONGO_DATABASE"
        value = "go2hal"
      }
      env {
        name = "ERROR_GROUP"
        value = "${var.error_group}"
      }
      env {
        "name" = "MONGO_SERVERS"
        value_from {
          config_map_key_ref {
            name = "mongo"
            key = "MONGO_SERVERS"
          }
        }
      }
      env {
        "name" = "MONGO_USER"
        value_from {
          config_map_key_ref {
            name = "mongo"
            key = "MONGO_USER"
          }
        }
      }
      env {
        "name" = "MONGO_REPLICA_SET"
        value_from {
          config_map_key_ref {
            name = "mongo"
            key = "MONGO_REPLICA_SET"
          }
        }
      }
      env {
        "name" = "MONGO_AUTH_SOURCE"
        value_from {
          config_map_key_ref {
            name = "mongo"
            key = "MONGO_AUTH_SOURCE"
          }
        }
      }
      env {
        "name" = "MONGO_SSL"
        value_from {
          config_map_key_ref {
            name = "mongo"
            key = "MONGO_SSL"
          }
        }
      }
      env{
        name = "MONGO_PASSWORD"
        value_from {
          secret_key_ref {
            name = "mongo"
            key = "MONGO_PASSWORD"
          }
        }
      }
      env {
        name = "HTTP_PROXY"
        value_from {
          secret_key_ref {
            name = "go2hal"
            key = "HTTP_PROXY"
          }
        }
      }
      env {
        name = "HTTPS_PROXY"
        value_from {
          secret_key_ref {
            name = "go2hal"
            key = "HTTPS_PROXY"
          }
        }
      }
      env {
        name = "BOT_KEY"
        value_from {
          secret_key_ref {
            name = "go2hal"
            key = "BOT_KEY"
          }
        }
      }
      security_context {
        privileged = "false"
        read_only_root_filesystem = "false"
      }
    }
    restart_policy = "Always"
  }

}

resource "kubernetes_service" "go2hal" {
  "metadata" {
    name = "go2hal"
    namespace = "go2hal"
  }
  "spec" {
    selector {
      app = "hal"
    }
    port {
      port = 80
      name = "http"
      target_port = 8000
      protocol = "TCP"
    }
  }
}

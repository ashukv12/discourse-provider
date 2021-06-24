terraform {
  required_providers {
    discourse = {
      version = "1"
      source  = "discourse.org/user/discourse"
    }
  }
}

provider "discourse" { 
  api_key="7c27290e2ccca7ae4427dfe518f68fb0659cd7df9908ff1cefbf034f9900a568"
  api_username="ashwinigaddagiwork"
  base_url="https://clevertaptest.trydiscourse.com"
}

resource "discourse_user" "user1" {
  email = "ashwinigaddagiwork@gmail.com"
  admin= "true"
}
/*
resource "discourse_user" "user1" {
  email = "ashutoshkverma12@gmail.com"
  name = "Ashutosh Verma"
  active = false
}
*/
output "user"{
    value = discourse_user.user1
}
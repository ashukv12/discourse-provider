This terraform provider enables create, read, update, delete, and import operations for discourse users.

## Requirements

* [Go](https://golang.org/doc/install) >= 1.16 (To build the provider plugin) <br>
* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x <br/>
* Application: [Discourse](https://www.discourse.org/pricing) (API is supported in Standard and Business plans) <br/>
* [Discourse API Documentation](https://docs.discourse.org/) 

## Application Account

### Setup 

1. Create a discourse account with your required subscription [Standard Plan/Business Account](https://www.discourse.org/pricing)<br>
2. Go to `Dashboard -> API -> API Key`. For our purpose we need to create an API Key. <br>

This app will provide us with the API Key which will be needed to configure our provider and make request. <br>
 
## API Authentication

1. authenticate API, we need a pair of credentials:  api_key,api_username, se_url, ifyUserID and expensifyUserSecret.
For this, go to https://www.expensify.com/tools/integrations/ and generate the credentials.
A pair of credentials: expensifyUserID and expensifyUserSecret will be generated and shown on the page.

## Building The Provider

1. Clone the repository, add all the dependencies and create a vendor directory that contains all dependencies. For this, run the following commands: 
 ```golang
cd terraform-provider-discourse
go mod init terraform-provider-discourse
```

2. Run `go mod vendor` to create a vendor directory that contains all the provider's dependencies. <br>

## Managing terraform plugins

For Linux:

1. Run the following command to create a vendor subdirectory which will comprise of  all provider dependencies. <br>
    ```
    ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
    ``` 
    Command: 
    ```bash
    export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
    mkdir -p ~/.terraform.d/plugins/discourse.org/user/discourse/1.0/$OS_ARCH
    ```

2. Run `go build -o terraform-provider-discourse`. This will save the binary file in the main/root directory. <br>

3. Run this command to move this binary file to appropriate location. <br>
```
  mv terraform-provider-discourse ~/.terraform.d/plugins/discourse.org/user/discourse/1.0/$OS_ARCH
```    
 <p align="center">
 [OR]
 </p><br>

3. Manually move the file from current directory to destination directory.<br>
 


## Working with terraform

### Application Credential Integration in terraform

1. Add terraform block and provider block as shown in example usage.
2. Get a pair of credentials: api_key, api_username, and base_url. For this, visit https://docs.discourse.org/.
3. Assign the above credentials to the respective field in the provider block.

### Basic Terraform Commands
1. `terraform init` - To initialize a working directory containing Terraform configuration files.
2. `terraform plan` - To create an execution plan. Displays the changes to be done.
3. `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the chages.

#### Create User
1. Add the user `email`, `name`, and `active`  in the respective field in resource block as shown in example usage.
2. Initialize the terraform provider `terraform init`
3. Check the changes applicable using `terraform plan` and apply using `terraform apply`
4. You will see that an account activation mail has been sent to the user.
5. Activate the account using the link provided in the mail, and you will see that a user has been successfully created.

#### Update the user
  User is allowed to update `name`, and `active`. Update the data, i.e., `name` or `active` of the user in the in the resource block of `main.tf` file and apply using `terraform apply`

#### Read the User Data
Add `data` and `output` blocks in the `main.tf` file and run `terraform plan` to read user data

#### Activate/Deactivate the user
Change the `active` field value from `false` to deactivate and `true` to activate and run `terraform apply`.

#### Delete the user
Delete the `resource` block of the particular user from `main.tf` file and run `terraform apply`.

#### Import a User Data
1. Write manually a resource configuration block for the User in `main.tf`, to which the imported object will be mapped.
2. Run the command `terraform import discourse_user.user1 [EMAIL_ID]`
3. Check for the attributes in the `.tfstate` file and fill them accordingly in resource block.


## Example Usage
```terraform
terraform {
  required_providers {
    discourse = {
      version = "1"
      source  = "discourse.org/user/discourse"
    }
  }
}

provider "discourse" {
  api_key = "_REPLACE_DISCOURSE_API_KEY_"
  api_username = "_REPLACE_DISCOURSE_API_USERNAME_"
  base_url = "_REPLACE_DISCOURSE_BASE_URL_"
}

resource "discourse_user" "user1" {
   email = "employee@domain.com"
   name = "Employee"
   active = true
}

data "discourse_user" "user1" {
  email = "employee@domain.com"
}

output "user1" {
  value = data.discourse_user.user1
}
```

## Undocumented APIs being used here:

#### Update a user by username

`https://{defaultHost}/u/{username}.json`

#### Activate the user

`http://{defaultHost}/admin/users/{USER_ID}/activate.json`

#### Get email id by using username 

`https://{defaultHost}/u/{username}/emails.json`


## Argument Reference:

* `email`       (Required, String)  - The email address of the user.
* `name`           (Optional, String)  - Name of the user in Discourse. 
* `active`         (Optional, Boolean) - If set to false, the user will be deactivated.

# **Test Forum API**
#### Test task of REST API

## **Tutorial**

- [Usage](#usage)
  - [Settings](#settings)
- [API](#api-points)
- [Authentication and Authorization](#authentication-and-authorization)

## **Usage**
#### Use ```go run``` or ```go build``` for launch application
#### For close application enter ```'server shutdown'``` to terminal

### **Settings**
#### Application settings represented by config.env file
```ini
HOST_ADDRESS=localhost:80   # TCP address for the server to listen on. Default: localhost:80
HASH_KEY=provider           # Key which hashing all tokens. Default: provider
USER_DB=utest               # database user
PASS_DB=12345               # database password
HOST_DB=localhost           # database host-address
PORT_DB=3306                # database host-port
NAME_DB=edudb               # database name
#Application Settings for Facebook Sign-In
FBA_CLIENT_ID=              # facebook application client ID
FBA_CLIENT_SECRET=          # facebook application client secret
    # If one of two previous fields is empty, then this type of authentication is not available 
FBA_REDIRECT_URL=http://localhost:80/auth/callback/facebook     #callbackURL
FBA_SCOPES=public_profile,email
FBA_AUTH_URL=https://www.facebook.com/v10.0/dialog/oauth
FBA_TOKEN_URL=https://graph.facebook.com/v10.0/oauth/access_token
FBA_API_VERSION=v10.0
#Application Settings for Google Sign-In
GA_CLIENT_ID=               # google application client ID
GA_CLIENT_SECRET=           # google application client secret
    # If one of two previous fields is empty, then this type of authentication is not available 
GA_REDIRECT_URL=http://localhost/auth/callback/google         #callbackURL
GA_SCOPES=https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/userinfo.profile,openid
GA_AUTH_URL=https://accounts.google.com/o/oauth2/auth
GA_TOKEN_URL=https://oauth2.googleapis.com/token
#Application Settings for Twitter Sign-In
TA_TWITTER_API_KEY=         # twitter application ID
TA_TWITTER_API_SECRET=      # twitter application secret
TA_TWITTER_TOKEN_KEY=       # twitter application token ID
TA_TWITTER_TOKEN_SECRET=    # twitter application token secret
    # If one of four previous fields is empty, then this type of authentication is not available 
TA_REDIRECT_URL=http://localhost:80/auth/callback/twitter     #callbackURL
TA_REQUEST_TOKEN_URL=https://api.twitter.com/oauth/request_token
TA_AUTH_URL=https://api.twitter.com/oauth/authenticate?oauth_token
TA_TOKEN_URL=https://api.twitter.com/oauth/access_token
```

## **API points**
#### Available api points, methods and query parameters -- APIPoint[method]
* [/posts [**GET**]](#posts-get)  
  _Available query parameters_:
  * userId
  * xml
* [/posts [**POST**]](#posts-post)
* [/posts [**PUT**] ](#posts-put)
* [/posts/#id [**GET**]](#posts-id-get)  
  _Available query parameters_:  
  * xml
* [/posts/#id [**DELETE**]](#posts-id-delete)
* [/posts/#id/comments [**GET**]](#posts-id-comments-get)  
  _Available query parameters_:  
  * xml
______________________________
* [/comments [**GET**]](#comments-get)  
  _Available query parameters_:
  * postId
  * xml
* [/comments [**POST**]](#comments-post)
* [/comments [**PUT**]](#comments-put)
* [/comments/#id [**GET**]](#comments-id-get)  
  _Available query parameters_:  
  * xml
* [/comments/#id [**DELETE**]](#comments-id-delete)
_______________________
* [/getapikey [**GET**]](#getapikey)  

### **Posts GET**
  List of all posts.   
  Available query parameters:  
  * userId (`/posts?userId=#id`) - list of posts, created by user with the given id (id is number)
  * xml (`/posts?xml`) - if xml parameter accepted result returned in xml format ; by default the list is returned in json format
### **Posts POST**
  Create post (requires authorization).  
  Parameters must be passed in the request body in json format:   
  ```json
  {
    "title" : "title",    
    "body": "body"
  }
  ```
  Both of parameters required.
### **Posts PUT**
  Updates an existing post (requires authorization).   
  Parameters must be passed in the request body in json format:   
  ```json
  {
    "id": id,   
    "title" : "title",
    "body": "body"
  }
  ```
  id - required. One of two: title or body required.
### **Posts ID GET**
  Get post by ID.   
  Available query parameters:    
  * xml (`/posts/#id?xml`) - if xml parameter accepted result returned in xml format ; by default the list is returned in json format
### **Posts ID DELETE**
  Delete post by ID (requires authorization).
### **Posts ID Comments GET**
  List of comments belonging to the post with the given ID.   
  Same as request [/comments?postId=#id [**GET**]](#comments-get)    
  Available query parameters:    
  * xml (`/posts/#id/comments?xml`) - if xml parameter accepted result returned in xml format ; by default the list is returned in json format
### **Comments GET**
  List of all comments.   
  Available query parameters:  
  * postId (`/comments?postId=#id`) - list of comments related to the post with the given id (id is number)
  * xml (`/comments?xml`) - if xml parameter accepted result returned in xml format ; by default the list is returned in json format
### **Comments POST**
  Create comment (requires authorization).  
  Parameters must be passed in the request body in json format:   
  ```json
  {
    "name" : "name",    
    "email": "email",
    "body": "body",
    "postId": postID
  }
  ```
  All parameters are required.
### **Comments PUT**
  Updates an existing comment (requires authorization).   
  Parameters must be passed in the request body in json format:   
  ```json
  {
    "id": id,   
    "name" : "name",
    "email": "email",
    "body": "body"
  }
  ```
  id - required. One of three: name, email or body required.
### **Comments ID GET**
  Get comment by ID.   
  Available query parameters:    
  * xml (`/comments/#id?xml`) - if xml parameter accepted result returned in xml format ; by default the list is returned in json format
### **Comments ID DELETE**
  Delete comment by ID (requires authorization).
### **GetAPIKey**
  Generate API Key for access without authentication (requires authorization). 
## **Authentication and Authorization**
  Aplication provide OAuth authentication through social media networks: Google `/auth/google`, Facebook `/auth/facebook`, Twitter `/auth/twitter`.  
  Authorization provide by:
  * `UAAT` cookie after authentication
  * `APIKey` HTTP header. API key can be generated `/getapikey` after authentication and used further without authentication until a new key is generated.
## **Licenses**
All source code is licensed under the [GNU License](https://github.com/ramesses-edu/nx_trainee_forum/blob/main/LICENSE)

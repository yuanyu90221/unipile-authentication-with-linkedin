#!/bin/bash

fetch_linked_user() {
  local token=$1
  local domain=$2
  `curl --request GET --url "$domain/api/v1/accounts" --header "X-API-KEY:$token" --header 'accept: application/json' | jq >> unipile/sample_response/response.json`
}

fetch_user_with_credential() {
  local username=$1
  local password=$2
  local domain=$3
  local token=$4
  `curl --request POST --url "$domain/api/v1/accounts" --header "X-API-KEY:$token" --header 'accept: application/json' --header 'content-type: application/json'  --data "{\"provider\": \"LINKEDIN\", \"username\": \"$username\", \"password\": \"$password\"}" | jq >> unipile/sample_response/login_response_v1.json`
}

fetch_checkout_otp() {
  local account_id=$1
  local otp=$2
  local domain=$3
  # local token=$4
  `curl --request POST --url "$domain/api/v1/accounts/checkpoint"  --header 'accept: application/json' --header 'content-type: application/json'  --data "{\"provider\": \"LINKEDIN\", \"code\": \"$otp\", \"account_id\": \"$account_id\"}" | jq >> unipile/sample_response/checkpoint_response.json`
}

fetch_user_with_cookie() {
  local access_token=$1
  local domain=$2
  local token=$3
  `curl --request POST --url "$domain/api/v1/accounts" --header "X-API-KEY:$token" --header 'accept: application/json' --header 'content-type: application/json'  --data "{\"provider\": \"LINKEDIN\", \"access_token\": \"$access_token\"}" | jq >> unipile/sample_response/cookie_response.json`
}

fetch_linked_user_with_user_id() {
  local account_id=$1
  local domain=$2
  local token=$3
  `curl --request GET --url "$domain/api/v1/accounts/$account_id" --header "X-API-KEY:$token" --header 'accept: application/json' | jq >> unipile/sample_response/fetch_user_response.json`
}
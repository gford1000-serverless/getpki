[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)


# GETPKI

This project implements two AWS Lambda function, along with associated AWS Gateway API supporting two methods:
* Return the API's RSA public key
* Return unique RSA key pairs on demand, signed using the API's RSA private key

Optionally the call can provide their own RSA public key in their request, and if provided then the RSA key pair will be encrypted using that key.

For more details on the underlying JSON structure and encryption protocol, see [PKIGEN](https://github.com/gford1000-go/pkigen).

The lambda functions are implemented using go 1.18, with AWS CloudFormation
script to create all the necessary AWS resources.

The CloudFormation script requires two parameters:

* The name of a valid AWS KMS key 
* A base64 encoded string, of a structured JSON object containing an RSA public/private key pair, encrypted with the above AWS KMS key

The lambda will require access to the KMS key, to decrypte the JSON object, as it starts.  The CloudFormation will create the necessary IAM Role and Policy to allow this.




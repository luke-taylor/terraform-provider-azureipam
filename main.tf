terraform {
  required_providers {
    azureipam = {
      source = "hashicorp.com/edu/azureipam"
    }
  }
}

provider "azureipam" {
  host = "https://ipam-xpmctiprtdfam.azurewebsites.net"
  token = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IktRMnRBY3JFN2xCYVZWR0JtYzVGb2JnZEpvNCJ9.eyJhdWQiOiI1NzY5NmQxMi0wZGUzLTQ2YTQtOGU1ZC0xY2NmMTgwNzY0YzAiLCJpc3MiOiJodHRwczovL2xvZ2luLm1pY3Jvc29mdG9ubGluZS5jb20vYjA3NmNjMWYtNmQ3ZS00MTAzLTk1ZjMtMGZmMGI4MTA4M2Q0L3YyLjAiLCJpYXQiOjE3MjMzMDk3MjgsIm5iZiI6MTcyMzMwOTcyOCwiZXhwIjoxNzIzMzEzNzE4LCJhaW8iOiJBWFFBaS84WEFBQUF4SEJoaGRxempNdStYb05MTFVXanBIVUFnMUF4YmJWWURpUnIyOStCVXk0TldPNEl2eUpVbU52VEs4b0lYM25JT3lCdDRFRE01Q3FWN3BicmRLQ084ZjdkK0Vyd3ZjOTNsOU9aVVNvTE51UzdUOWxpVXh5SXFsZG5kSGE5T1NwK2dlclNiempvSWhRUXVzZ0ppazhxVGc9PSIsImF6cCI6IjI3OTEzNThjLWU2MDYtNGVlZC1hYjczLWU4NDEyZjNlOTJhYyIsImF6cGFjciI6IjAiLCJuYW1lIjoiTHVrZSBUYXlsb3IiLCJvaWQiOiJiMjc1MDNhOS02ZDE1LTQ5YTgtODlkYi1iYTM5MzRjYmI1MDciLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJsdWtlQE1uZ0Vudk1DQVA1NzMwMjUub25taWNyb3NvZnQuY29tIiwicmgiOiIwLkFYd0FIOHgyc0g1dEEwR1Y4d193dUJDRDFCSnRhVmZqRGFSR2psMGN6eGdIWk1DN0FCRS4iLCJzY3AiOiJhY2Nlc3NfYXNfdXNlciIsInN1YiI6IlNUMTl1NUduSnloMFJNVi02ZFUyb2wxNW5nYURibHdLaW9PdmVYMm9vSEUiLCJ0aWQiOiJiMDc2Y2MxZi02ZDdlLTQxMDMtOTVmMy0wZmYwYjgxMDgzZDQiLCJ1dGkiOiJyR3BXUWtoRjRVRzVLUXQ5d1dNX0FBIiwidmVyIjoiMi4wIn0.MFPzSI-rplc5bjcY1jSOC2FbI4gK1mbIASCWI-Z46K5KfftDcSo1xJy2QQyFwbXu9QOjuBpCoCKG-T8X1qO5aLv6GTtM3Z6JWCf565mxVxv-4eNo7TKupYkhnM3f20QNc5WF4gHsAxBLRRNKXiIXh5LrcYp0VCW0A3N39u_nTBLfKXnZUP9IbEyQtCwgNmpXhxmhhRdXmKHoM-rtYPvxm209f0gkJ6pAvK7GN4mxpixMeKp4CtxK7yXCE3ccpkHCDdEgTNM2h2RUuhpas1OJ3x7O7oZv6KwVVdkDduqDwm6x1kDXiBdPRNfgl9PXP6oJ_UoJpXnWfjS19tbDoSkdqg"
}


data "azureipam_admins" "example" {
}

output "names" {
  value = data.azureipam_admins.example
  
}
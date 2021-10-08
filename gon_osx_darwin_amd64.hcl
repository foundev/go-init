source = ["./dist/go-init-osx-amd_darwin_amd64/go-init"]
bundle_id = "pro.foundev.goinit"

sign {
  application_identity = "Developer ID Application: Ryan Svihla (8FLL83XJM2)"
}

zip {
  output_path = "./dist/go-init-osx-amd64-signed.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

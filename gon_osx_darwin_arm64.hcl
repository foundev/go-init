source = ["./dist/go-init-osx-arm_darwin_arm64/go-init"]
bundle_id = "pro.foundev.goinit"

sign {
  application_identity = "DDF52A8D387B7E77F18F86043FF2AC7AED277179"
}

zip {
  output_path = "./dist/go-init-osx-arm-signed.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

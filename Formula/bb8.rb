class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.10/bb8_0.1.10_macOS_amd64.tar.gz"
  version "0.1.10"
  sha256 "3c84c3458653c8071affd88b0e8a8338efb9826c4d2fba387951e75140397aa8"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

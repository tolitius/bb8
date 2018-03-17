class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.14/bb8_0.1.14_macOS_amd64.tar.gz"
  version "0.1.14"
  sha256 "e05434469e65f6aca63ced7189e07a31057d02c42435263b85976f3dd5f0eb35"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

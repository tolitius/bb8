class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.18/bb8_0.1.18_macOS_amd64.tar.gz"
  version "0.1.18"
  sha256 "86ee50e44ba7d665595aba2c9e443b2c15429f7e2c0bb897443d24f11d0b1a9e"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

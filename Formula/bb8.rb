class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.13/bb8_0.1.13_macOS_amd64.tar.gz"
  version "0.1.13"
  sha256 "4d2301784e26c41011411badb66d872e5c49f2bf229599392868650a36b579eb"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

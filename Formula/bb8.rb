class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage ""
  url "https://github.com/tolitius/bb8/releases/download/v0.1.3/bb8_0.1.3_macOS_amd64.tar.gz"
  version "0.1.3"
  sha256 "a760f2ba87346a28f05c943f019d40e465f7247e3b3a96832cd627d8d232d9ea"

  def install
    bin.install "bb"
  end
end

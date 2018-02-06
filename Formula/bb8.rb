class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.7/bb8_0.1.7_macOS_amd64.tar.gz"
  version "0.1.7"
  sha256 "ca53d3dc63cdde39d2120f5d8718d803754f48e7088a611e4a5448b51900330a"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.16/bb8_0.1.16_macOS_amd64.tar.gz"
  version "0.1.16"
  sha256 "9c58586b3b645388488d11c12599a5b91ed76a49ebddc9e35699cc7450ea51d6"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

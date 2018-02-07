class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.8/bb8_0.1.8_macOS_amd64.tar.gz"
  version "0.1.8"
  sha256 "7cb71a56c48f6025081f61022316c1284b36670e5f04f209448248d140d99308"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

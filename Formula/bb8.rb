class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.9/bb8_0.1.9_macOS_amd64.tar.gz"
  version "0.1.9"
  sha256 "91a06c62f145b6bf879759d89844cdfcf8fdd539df89b9caf82599c87ec56b9f"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

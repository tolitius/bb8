class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.17/bb8_0.1.17_macOS_amd64.tar.gz"
  version "0.1.17"
  sha256 "bd3107c697324e5ca8e5c92684b52d114939e0b40ec7faad30aa5aee3e6e5463"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

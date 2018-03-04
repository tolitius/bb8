class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.12/bb8_0.1.12_macOS_amd64.tar.gz"
  version "0.1.12"
  sha256 "4e8924d3d2100dd24c3d91208bc2dc84db9f83fffc488aa2f7bf9f9d4865443a"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

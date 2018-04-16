class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.15/bb8_0.1.15_macOS_amd64.tar.gz"
  version "0.1.15"
  sha256 "2a20151e543954c91015d93adbbba4f44e1575e29af27d6be531a9753d59057d"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

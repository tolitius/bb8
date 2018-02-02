class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage ""
  url "https://github.com/tolitius/bb8/releases/download/v0.1.4/bb8_0.1.4_macOS_amd64.tar.gz"
  version "0.1.4"
  sha256 "a2d8dba9199f533f80f3aa323d726406a662af7b568a1e9aadd051575d6fa741"

  def install
    bin.install "bb"
  end
end

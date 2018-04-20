class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.18/bb8_0.1.18_macOS_amd64.tar.gz"
  version "0.1.18"
  sha256 "5b6c1dfbeceecf59811196a3c92582cfb7268e9857eede83a2675632b333a79e"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

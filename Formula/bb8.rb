class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.6/bb8_0.1.6_macOS_amd64.tar.gz"
  version "0.1.6"
  sha256 "d09c537e81d2258b731fa56422e5c4199e20e127d2be6172b6aadf5b4fbefdb9"

  def install
    bin.install "bb"
  end

  test do
    
  end
end

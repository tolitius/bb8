class Bb8 < Formula
  desc "a command line interface to Stellar networks"
  homepage "https://github.com/tolitius/bb8"
  url "https://github.com/tolitius/bb8/releases/download/v0.1.5/bb8_0.1.5_macOS_amd64.tar.gz"
  version "0.1.5"
  sha256 "50b71ce4c3383c9c3c286cf4de4fe71d9a23955fd3811166ebcd2e6ac543505a"

  def install
    bin.install "bb"
  end

  test do
    bb version
  end
end

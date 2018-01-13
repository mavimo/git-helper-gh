class GitHelperGh < Formula
  desc "Git helper for GitHub"
  homepage "https://github.com/mavimo/git-helper-gh"
  url "https://github.com/mavimo/git-helper-gh/archive/master.zip"
  version "0.0.1"
  # sha256 ""

  def install
    bin.install "bin/git-gh-start"
    bin.install "bin/git-gh-pr"
  end

  test do
    system "false"
  end
end

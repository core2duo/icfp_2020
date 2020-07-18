#pragma once

#include <string>
#include <vector>
#include <memory>

#include "tokens.h"
#include "appliers/cons.h"

std::unique_ptr<Token> tokenize(const std::string& s) {
  if (s == "nil") return std::make_unique<NilToken>();

  if (s == "ap") return std::make_unique<ApplicationToken>();

  if (s == "cons") return std::make_unique<ConsFunctionToken>();

  if (s == "=") return std::make_unique<AssignmentToken>();
  if (!s.empty() && s[0] == ':') return std::make_unique<VariableToken>(s);

  return std::make_unique<IntegerToken>(std::stoll(s));
}

std::vector<std::unique_ptr<Token>> tokenize(const std::vector<std::string>& v) {
  std::vector<std::unique_ptr<Token>> result;
  result.reserve(v.size());
  for (const auto& s : v) {
    std::cout << s << std::endl;
    result.push_back(std::move(tokenize(s)));
  }
  return result;
}

std::vector<std::string> split(const std::string& s) {
  std::vector<std::string> result;
  std::string buffer;

  if (s.empty()) {
    return result;
  }

  for (auto c : s) {
    if (c == ' ') {
      result.push_back(std::move(buffer));
      buffer.clear();
    } else {
      buffer += c;
    }
  }
  result.push_back(std::move(buffer));

  return result;
}
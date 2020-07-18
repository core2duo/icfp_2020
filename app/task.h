#pragma once

#include <memory>
#include <string>
#include <vector>
#include <iostream>
#include <unordered_map>

#include "util.h"

std::unordered_map<size_t, Token> context;

void eval(const std::string& s) {
  auto tokens = tokenize(split(s));
  std::vector<std::unique_ptr<Token>> stack;
  stack.reserve(tokens.size());

  for (auto it = tokens.rbegin(); it != tokens.rend(); ++it) {
    print(it->get());

    if (auto ap = dynamic_cast<ApplicationToken*>(it->get())) {

      if (stack.size() < 2) {
        // Whoops, something went wrong
        std::cout << "stack size is less than 2" << std::endl;
        exit(1);
      }
      auto f = std::move(stack.back());
      stack.pop_back();
      auto op = std::move(stack.back());
      stack.pop_back();

      if (auto func = dynamic_cast<FunctionToken*>(f.get())) {
        auto next_f = func->apply(std::move(op));
        stack.push_back(std::move(next_f));
      }

    } else {
      stack.push_back(std::move(*it));
    }
  }

  std::cout << "Size " << stack.size() << std::endl;

  if (auto p = dynamic_cast<ConsFunctionToken*>(stack[2].get())) {
    std::cout << "here" << std::endl;
    for (const auto& t : p->list.value) {
      print(t.get());
    }
  }

}
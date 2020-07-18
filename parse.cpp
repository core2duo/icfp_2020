#include <algorithm>
#include <deque>
#include <fstream>
#include <functional>
#include <iostream>
#include <iterator>
#include <map>
#include <memory>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

using namespace std;

#ifndef NDEBUG
#   define ASSERT(condition, message) \
    do { \
        if (! (condition)) { \
            std::string err = "Assertion `" #condition "` failed with `" #message "` in line " + std::to_string(__LINE__); \
            throw std::runtime_error(err); \
        } \
    } while (false)
#else
#   define ASSERT(condition, message) do { } while (false)
#endif

// Ast things.
enum class Type {
  Fn = 0,
  Pair,
  Number,
  Variable,
  Nill,
  Any,
};

enum class Primitives {
  Ap = 0,
  Add,
  B,
  C,
  Car,
  Cdr,
  Cons,
  Div,
  Eq,
  Galaxy,
  I,
  Isnil,
  Lt,
  Mul,
  Neg,
  Nil,
  S,
  T,
  F,
  Number,
  Variable,
};

using Fn = Primitives;

constexpr const char* FnToStr(Fn fn) noexcept {
  switch(fn) {
    case Fn::Ap: return "ap";
    case Fn::Add: return "add";
    case Fn::B: return "b";
    case Fn::C: return "c";
    case Fn::Car: return "car";
    case Fn::Cdr: return "cdr";
    case Fn::Cons: return "cons";
    case Fn::Div: return "div";
    case Fn::Eq: return "eq";
    case Fn::Galaxy: return "galaxy";
    case Fn::I: return "i";
    case Fn::Isnil: return "isnil";
    case Fn::Lt: return "lt";
    case Fn::Mul: return "mul";
    case Fn::Neg: return "neg";
    case Fn::Nil: return "nil";
    case Fn::S: return "s";
    case Fn::T: return "t";
    case Fn::F: return "f";
    case Fn::Number: return "Number";
    case Fn::Variable: return "Variable";
  }
}

Fn StrToFn(const string& str) {
  if (str == "ap") return Fn::Ap;
  if (str == "add") return Fn::Add;
  if (str == "b") return Fn::B;
  if (str == "c") return Fn::C;
  if (str == "car") return Fn::Car;
  if (str == "cdr") return Fn::Cdr;
  if (str == "cons") return Fn::Cons;
  if (str == "div") return Fn::Div;
  if (str == "eq") return Fn::Eq;
  if (str == "galaxy") return Fn::Galaxy;
  if (str == "i") return Fn::I;
  if (str == "isnil") return Fn::Isnil;
  if (str == "lt") return Fn::Lt;
  if (str == "mul") return Fn::Mul;
  if (str == "neg") return Fn::Neg;
  if (str == "nil") return Fn::Nil;
  if (str == "s") return Fn::S;
  if (str == "t") return Fn::T;
  if (str == "f") return Fn::F;

  if (!str.empty() && str[0] == ':') {
    return Fn::Variable;
  }

  return Fn::Number;
}

struct Node;

struct Environment {
  map<int64_t, Node*> nodes;
};

using Env = Environment;

struct NodeNill;
struct NodeFn;
struct NodeNumber;
struct NodePair;
struct NodeVariable;

struct Node {
  virtual ~Node() = default;
  virtual Type GetType() const = 0;
  virtual NodeNill* AsNill() { return nullptr; }
  virtual NodeFn* AsFn() { return nullptr; }
  virtual NodeNumber* AsNumber() { return nullptr; }
  virtual NodePair* AsPair() { return nullptr; }
  virtual NodeVariable* AsVariable() { return nullptr; }
};

struct NodeNill : public Node {
  Type GetType() const override { return Type::Nill; }
  NodeNill* AsNill() override { return this; }
};

NodeNill* NILL() {
  static NodeNill nill;
  return &nill;
}

struct NodeFn: public Node {
  Type GetType() const override { return Type::Fn; }

  NodeFn* AsFn() override { return this; }

  Fn fn;
  vector<Node*> args;
};

NodeFn* TRUE() {
  static NodeFn node_true;
  node_true.fn = Fn::T;
  return &node_true;
}

NodeFn* FALSE() {
  static NodeFn node_true;
  node_true.fn = Fn::F;
  return &node_true;
}

struct NodeNumber : public Node {
  Type GetType() const override { return Type::Number; }
  NodeNumber* AsNumber() override { return this; }

  int64_t number;
};

struct NodePair : public Node {
  Type GetType() const override { return Type::Pair; }
  NodePair* AsPair() override { return this; }

  Node* left = nullptr;
  Node* right = nullptr;
};

struct NodeVariable : public Node {
  NodeVariable * AsVariable() override { return this; }
  Type GetType() const override { return Type::Variable; }

  int64_t id;
};

using Block = vector<string>;

template <typename T>
struct Visitor {
  T Visit(Node* node) {
    ASSERT(node != nullptr, "must be real");
    switch(node->GetType()) {
      case Type::Nill:
        return VisitNill(node->AsNill()); 
      case Type::Fn:
        return VisitFn(node->AsFn());
      case Type::Number:
        return VisitNumber(node->AsNumber());
      case Type::Pair:
        return VisitPair(node->AsPair());
      case Type::Variable:
        return VisitVariable(node->AsVariable());
    }
  }

  virtual T VisitNill(NodeNill* node) = 0;
  virtual T VisitFn(NodeFn* node) = 0;
  virtual T VisitNumber(NodeNumber* node) = 0;
  virtual T VisitPair(NodePair* node) = 0;
  virtual T VisitVariable(NodeVariable* node) = 0;

  virtual ~Visitor() = default;
};

string Center(int width, const string& str) {
  int len = str.length();
  if (width < len) {
    return str;
  }

  int diff = width - len;
  int pad1 = diff/2;
  int pad2 = diff - pad1;
  return string(pad1, ' ') + str + string(pad2, ' ');
}

struct PrintVisitor : public Visitor<Block> {
  Block VisitNill(NodeNill* node) override {
    return {"NILL"};
  }

  Block Compose(const string& top, vector<Node*> nodes) {
    vector<Block> blocks; 

    size_t width = 0;
    size_t height = 0;

    size_t left_pad = 0;
    for (size_t i = 0; i < nodes.size(); ++i) {
      Block block = Visit(nodes[i]);
      ASSERT(!block.empty(), "mustn't be empty");
      width += block[0].size();
      height = max(height, block.size());
      if (i == nodes.size() / 2 - 1) {
        left_pad = width;
      }
      blocks.push_back(std::move(block));
    }
    size_t right_pad = width - left_pad;
    width += top.size() + 2;

    Block result; 
    stringstream ss1;
    ss1 << string(left_pad, ' ');
    ss1 << " ";
    ss1 << top;
    ss1 << " ";
    ss1 << string(right_pad, ' ');
    result.push_back(ss1.str());
    ASSERT(result.back().size() == width, "must match");

    for (size_t i = 0; i < height; ++i) {
      stringstream ss;

      for (size_t j = 0; j < blocks.size(); ++j) {
        auto& block = blocks[j];

        if (i < block.size()) { 
          ss << block[i];
        } else {
          ss << string(block[0].size(), ' ');
        }

        if (j == blocks.size() / 2 - 1) {
          ss << string(top.size() + 2, ' ');
        }
      }

      result.push_back(ss.str());
      ASSERT(result.back().size() == width, "must be aligned");
    }

    return result;
  }

  Block VisitFn(NodeFn* node) override {
    return Compose(FnToStr(node->fn), node->args);
  }

  Block VisitNumber(NodeNumber* node) {
    return {to_string(node->number)};
  }

  Block VisitPair(NodePair* node) {
    return Compose("pair", {node->left, node->right});
  }

  Block VisitVariable(NodeVariable* node) {
    return {":" + to_string(node->id)};
  }

  string Print(Node* node) {
    Block block = Visit(node); 
    stringstream ss;
    for (string& line : block) {
      ss << line << "\n";
    }
    return ss.str();
  }
};

Node* Parse(deque<string>& tokens);

Node* ParseFn(Fn fn, size_t nb_args, deque<string>& tokens) {
  auto node = new NodeFn{};
  node->fn = fn; 
  
  for (size_t i = 0; i < nb_args; ++i) {
    auto arg = Parse(tokens);
    if (arg == nullptr) {
      return node;
    }
    node->args.push_back(arg);
  }

  return node;
}

Node* Parse(deque<string>& tokens) {
  if (tokens.empty()) {
    return nullptr;
  }

  string str = std::move(tokens.front());
  tokens.pop_front();

  struct Entry {
    Fn fn;
    size_t args;
  };

  static const map<string, Entry> fns = {
    {"add", {Fn::Add, 0}},
    {"ap", {Fn::Ap, 2}},
    {"b", {Fn::B, 0}},
    {"c", {Fn::C, 0}},
    {"car", {Fn::Car, 0}},
    {"cdr", {Fn::Cdr, 0}},
    {"cons", {Fn::Cons, 0}},
    {"div", {Fn::Div, 0}},
    {"eq", {Fn::Eq, 0}},
    {"i", {Fn::I, 0}},
    {"isnil", {Fn::Isnil, 0}},
    {"lt", {Fn::Lt, 0}},
    {"mul", {Fn::Mul, 0}},
    {"neg", {Fn::Neg, 0}},
    {"s", {Fn::S, 0}},
    {"t", {Fn::T, 0}},
  };

  if (auto it = fns.find(str); it != fns.end()) {
    return ParseFn(it->second.fn, it->second.args, tokens);
  }

  if (str == "nil") {
    return NILL();
  }

  if (!str.empty() && str[0] == ':') {
    auto var = new NodeVariable{};
    var->id = std::stoll(str.substr(1));
  }

  auto number = new NodeNumber{};
  number->number = std::stoll(str);
  return number;
}

struct EvalVisitor : Visitor<Node*> {
  Node* VisitNill(NodeNill* node) override {
    return NILL();
  }

  Node* VisitFn(NodeFn* node) override {
    struct Entry {
      vector<Type> types;
      function<Node*(EvalVisitor&, NodeFn*)> f;
    };

    static const map<Fn, Entry> entries = {
      {Fn::Add, {
        {Type::Number, Type::Number},
        [](EvalVisitor&, NodeFn* f) {
          NodeNumber* a0 = f->args[0]->AsNumber();
          NodeNumber* a1 = f->args[1]->AsNumber();
          NodeNumber* n = new NodeNumber; 
          n->number = a0->number + a1->number;
          delete a0;
          delete a1;
          delete f;
          return n;
        }
      }},
      {Fn::Mul, {
        {Type::Number, Type::Number},
        [](EvalVisitor&, NodeFn* f) {
          NodeNumber* a0 = f->args[0]->AsNumber();
          NodeNumber* a1 = f->args[1]->AsNumber();
          NodeNumber* n = new NodeNumber; 
          n->number = a0->number * a1->number;
          delete a0;
          delete a1;
          delete f;
          return n;
        }
      }},
      {Fn::Div, {
        {Type::Number, Type::Number},
        [](EvalVisitor&, NodeFn* f) {
          NodeNumber* a0 = f->args[0]->AsNumber();
          NodeNumber* a1 = f->args[1]->AsNumber();
          NodeNumber* n = new NodeNumber; 
          n->number = a0->number / a1->number;
          delete a0;
          delete a1;
          delete f;
          return n;
        }
      }},
      {Fn::Car, {
        {Type::Pair},
        [](EvalVisitor& v, NodeFn* f) {
          NodePair* a0 = f->args[0]->AsPair();
          Node* n = v.Visit(a0->left);
          delete f;
          // Leak Need gc.
          // delete a0->right
          return n;
        }
      }},
      {Fn::Cdr, {
        {Type::Pair},
        [](EvalVisitor& v, NodeFn* f) {
          NodePair* a0 = f->args[0]->AsPair();
          Node* n = v.Visit(a0->right);
          delete f;
          // Leak Need gc.
          // delete a0->left;
          return n;
        }
      }},
      {Fn::Eq, {
        {Type::Number, Type::Number},
        [](EvalVisitor&, NodeFn* f) {
          NodeNumber* a0 = f->args[0]->AsNumber();
          NodeNumber* a1 = f->args[1]->AsNumber();
          Node* n = (a0->number == a1->number) ? TRUE() : FALSE();
          delete a0;
          delete a1;
          delete f;
          return n;
        }
      }},
      {Fn::Lt, {
        {Type::Number, Type::Number},
        [](EvalVisitor&, NodeFn* f) {
          NodeNumber* a0 = f->args[0]->AsNumber();
          NodeNumber* a1 = f->args[1]->AsNumber();
          Node* n = (a0->number < a1->number) ? TRUE() : FALSE();
          delete a0;
          delete a1;
          delete f;
          return n;
        }
      }},
      {Fn::Isnil, {
        {Type::Any},
        [](EvalVisitor&, NodeFn* f) {
          Node* n = (f->args[0]->GetType() == Type::Nill) ? TRUE() : FALSE();
          delete f;
          return n;
        }
      }},
      {Fn::T, {
        {Type::Any, Type::Any},
        [](EvalVisitor&, NodeFn* f) {
          Node* a0 = f->args[0];
          Node* a1 = f->args[1];
          delete a1;
          delete f;
          return a0;
        }
      }},
      {Fn::F, {
        {Type::Any, Type::Any},
        [](EvalVisitor&, NodeFn* f) {
          Node* a0 = f->args[0];
          Node* a1 = f->args[1];
          delete a0;
          delete f;
          return a1;
        }
      }},
      {Fn::Ap, {
        {Type::Fn, Type::Any},
        [](EvalVisitor& v, NodeFn* f) {
          NodeFn* a0 = f->args[0]->AsFn();
          Node* a1 = f->args[1];
          a0->args.push_back(a1);
          Node* n = v.Visit(a0);
          delete f; 
          return n;
        }
      }},
      {Fn::Neg, {
        {Type::Number},
        [](EvalVisitor& v, NodeFn* f) {
          NodeNumber* n = f->args[0]->AsNumber();
          n->number = -n->number;
          delete f;
          return n;
        }
      }},
      {Fn::Cons, {
        {Type::Any, Type::Any},
        [](EvalVisitor& v, NodeFn* f) {
          NodePair* pair = new NodePair;
          pair->left = f->args[0];
          pair->right = f->args[1];
          delete f;
          return pair;
        }
      }},
      {Fn::B, {
        {Type::Fn, Type::Fn, Type::Any},
        [](EvalVisitor& v, NodeFn* f) {
          // b - x y z = x (y z)
          NodeFn* x = f->args[0]->AsFn();
          NodeFn* y = f->args[1]->AsFn();
          Node* z = f->args[2];
          y->args.push_back(z);
          x->args.push_back(v.Visit(y));
          Node* n = v.Visit(x);
          delete f;
          return n;
        }
      }},
      {Fn::C, {
        {Type::Fn, Type::Any, Type::Any},
        [](EvalVisitor& v, NodeFn* f) {
          // c - x y z = (x z) y
          NodeFn* x = f->args[0]->AsFn();
          Node* y = f->args[1];
          Node* z = f->args[2];
          x->args.push_back(z);
          Node* g = v.Visit(x);
          ASSERT(g->GetType() == Type::Fn, "Must be a function");
          g->AsFn()->args.push_back(y);
          Node* n = v.Visit(g);
          delete g;
          delete f;
          return n;
        }
      }},
      {Fn::S, {
        {Type::Fn, Type::Fn, Type::Any},
        [](EvalVisitor& v, NodeFn* f) {
          // s x y z = x(z)(y(z))
          NodeFn* x = f->args[0]->AsFn();
          NodeFn* y = f->args[1]->AsFn();
          Node* z = f->args[2];
          x->args.push_back(z);
          Node* g = v.Visit(x);
          ASSERT(g->GetType() == Type::Fn, "Must be a function");
          y->args.push_back(z);
          g->AsFn()->args.push_back(v.Visit(y));
          Node* n = v.Visit(g);
          delete g;
          delete f;
          return n;
        }
      }},
    };

    auto it = entries.find(node->fn);
    if (it == entries.end()) {
      cout << "Didn't found " << FnToStr(node->fn) << endl;
      return node;
    }

    for (size_t i = 0; i < node->args.size(); ++i) {
      node->args[i] = Visit(node->args[i]);
    }

    //cout << "f " << FnToStr(node->fn) << " ";
    //cout << "m " << node->args.size() << " " << it->second.types.size() << endl;
    ASSERT(node->args.size() <= it->second.types.size(), "wrong size");

    if (node->args.size() != it->second.types.size()) {
      return node;
    }

    for (size_t i = 0; i < node->args.size(); ++i) {
      if (it->second.types[i] == Type::Any) {
        continue;
      }

      if (it->second.types[i] != node->args[i]->GetType()) {
        return node;
      }
    }

    return it->second.f(*this, node);
  }

  Node* VisitNumber(NodeNumber* node) override {
    return node;
  }

  Node* VisitPair(NodePair* node) override {
    return node;
  }

  Node* VisitVariable(NodeVariable* node) override {
    return node;
  }
};


int main(int argc, char** argv) {
  if (argc < 2) { 
    return -1;
  }

  cout << "file " << argv[1] << endl;
  ifstream file(argv[1]);
  std::string str;
  vector<deque<string>> lines;
  cout << "here" << endl;
  while (std::getline(file, str)) {
    istringstream iss(str);
    lines.emplace_back();
    copy(istream_iterator<string>(iss),
         istream_iterator<string>(),
         back_inserter(lines.back()));
  }

  //str = "ap ap cons 7 ap ap cons 123229502148636 nil";
  //str = "ap ap cons 1 ap ap cons 2 ap ap cons 3 ap ap cons 4 nil";
  //str = "ap ap cons 2 ap ap cons 7 nil";
  str = "ap ap s ap ap b c isnil car";
  istringstream iss(str);
  deque<string> line;
  copy(istream_iterator<string>(iss),
       istream_iterator<string>(),
       back_inserter(line));
  for (auto& a : line) {
    cout << a << " ";
  }
  cout << endl;
  cout << endl;
  Node* node = Parse(line);
  ASSERT(line.empty(), "must be empty");
  ASSERT(node->GetType() == Type::Fn && node->AsFn()->fn == Fn::Ap, "must be ap");
  PrintVisitor print_visitor;
  cout << print_visitor.Print(node) << endl;
  cout << endl;
  cout << endl;

  NodeNumber* number = new NodeNumber;
  number->number = 5;
  EvalVisitor eval_visitor;
  cout << "eval start" << endl;
  Node* root = eval_visitor.Visit(node);
  cout << "eval end" << endl;
  cout << print_visitor.Print(root) << endl;
  //NodeNumber number(5);
  //NodeVariable variable(10);
  //NodePair pair(&number, &variable);
  //NodeFn root(Fn::Ap);
  //root.args = {
  //  NILL(),
  //  &number,
  //  &variable,
  //  &pair,
  //};

}

#include "crow.h"

int main() { // тут будет что-то нормальное, но пока это
    crow::SimpleApp app;

    CROW_ROUTE(app, "/")([](){
        return "kek";
    });

    app.port(18080).run();
}